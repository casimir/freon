import 'dart:convert';

import 'package:cadanse/cadanse.dart';
import 'package:cadanse/components/widgets/error.dart';
import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:json_annotation/json_annotation.dart';

import '../services/freon.dart';

part 'forms.g.dart';

class ResourceForm extends ConsumerStatefulWidget {
  const ResourceForm({
    super.key,
    required this.resourcePath,
    this.resourceKey,
    this.resourceSchema,
  });

  final String resourcePath;
  final String? resourceKey;
  final String? resourceSchema;

  String get baseUrl => '/control$resourcePath';
  String get resourceUrl =>
      '$baseUrl${resourceKey != null ? '/$resourceKey' : ''}';
  ObjectSchemaPath get osp =>
      ObjectSchemaPath(resourceUrl, resourceSchema ?? '$baseUrl/schema');

  @override
  ConsumerState<ResourceForm> createState() => _FormState();
}

class _FormState extends ConsumerState<ResourceForm> {
  final GlobalKey<FormBuilderState> _formKey = GlobalKey<FormBuilderState>();
  bool _isUpdating = false;
  bool _obscureText = true;

  @override
  Widget build(BuildContext context) {
    return ref.watch(jsonFetcher(widget.osp)).when(
          data: (data) => _buildForm(data),
          error: (error, _) => ErrorScreen(error: error as Exception),
          loading: () => const Center(child: CircularProgressIndicator()),
        );
  }

  Widget _buildForm(List<dynamic> data) {
    final fields =
        data.map((it) => FormFieldValue.fromJson(Map.from(it))).toList();
    final isCreation = fields.any((it) => it.name == 'ID' && it.value == 0);
    return Center(
      child: Column(
        children: [
          FormBuilder(
            key: _formKey,
            child: ListView(
              shrinkWrap: true,
              children: fields
                  .where((it) => !(isCreation && it.readonly))
                  .map((it) => _buildFieldEntry(it))
                  .toList(),
            ),
          ),
          C.spacers.verticalContent,
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Row(
                children: [
                  Switch(
                    value: !_obscureText,
                    onChanged: (value) => setState(() {
                      _obscureText = !value;
                    }),
                  ),
                  C.spacers.verticalContent,
                  const Text('Show sensitive data'),
                ],
              ),
              ElevatedButton(
                onPressed: () async {
                  try {
                    setState(() => _isUpdating = true);
                    if (_formKey.currentState!.saveAndValidate()) {
                      final data = _formKey.currentState!.value;
                      await jsonUpload(widget.resourceUrl, data);
                    }
                  } on FreonError catch (e) {
                    if (context.mounted) {
                      ScaffoldMessenger.of(context).showSnackBar(SnackBar(
                        content: Text(e.toString()),
                      ));
                    }
                  } finally {
                    _isUpdating = false;
                    ref.invalidate(jsonFetcher(widget.osp));
                  }
                },
                child: _isUpdating
                    ? const CircularProgressIndicator()
                    : Text(isCreation ? 'Create' : 'Update'),
              ),
            ],
          )
        ],
      ),
    );
  }

  Widget _buildFieldEntry(FormFieldValue field) {
    return FormBuilderTextField(
      name: field.name,
      decoration: InputDecoration(labelText: field.name),
      initialValue: field.value.toString(),
      readOnly: field.readonly,
      obscureText: field.obfuscate && _obscureText,
    );
  }
}

@JsonSerializable(createToJson: false)
class FormFieldValue {
  FormFieldValue(this.name, this.value, this.readonly, this.obfuscate);

  final String name;
  final dynamic value;
  final bool readonly;
  final bool obfuscate;

  factory FormFieldValue.fromJson(Map<String, dynamic> json) =>
      _$FormFieldValueFromJson(json);
}
