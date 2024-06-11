import 'package:cadanse/cadanse.dart';
import 'package:cadanse/components/widgets/error.dart';
import 'package:cadanse/layout.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../services/freon.dart';
import 'forms.dart';

class ResourceList extends ConsumerStatefulWidget {
  const ResourceList({
    super.key,
    required this.resourcePath,
    this.resourceSchema,
  });

  final String resourcePath;
  final String? resourceSchema;

  @override
  ConsumerState<ResourceList> createState() => _ResourceListState();

  String get baseUrl => '/control$resourcePath';
  String resourceUrl([String? resourceKey]) =>
      '$baseUrl${resourceKey != null ? '/$resourceKey' : ''}';
  ObjectSchemaPath osp([String? resourceKey]) => ObjectSchemaPath(
        resourceUrl(resourceKey),
        resourceSchema ?? '$baseUrl/schema',
      );
}

class _ResourceListState extends ConsumerState<ResourceList> {
  @override
  Widget build(BuildContext context) {
    return ref.watch(jsonFetcher(widget.osp())).when(
          data: (data) =>
              data.isNotEmpty ? _buildList(data) : _buildAddButton(context),
          error: (error, _) => ErrorScreen(error: error as Exception),
          loading: () => const Center(child: CircularProgressIndicator()),
        );
  }

  Widget _buildList(List<dynamic> data) {
    final items = data
        .map((e) => Map<String, FormFieldValue>.fromIterable(
              e.map((it) => FormFieldValue.fromJson(Map.from(it))),
              key: (it) => it.name,
            ))
        .toList();
    return Column(
      children: [
        Row(children: [const Spacer(), _buildAddButton(context)]),
        ListView.builder(
          shrinkWrap: true,
          itemCount: items.length,
          itemBuilder: (context, index) {
            final item = items[index];
            var name = item['Name']!.value;
            if (name == null || name.isEmpty) {
              name = '<No name>';
            }
            return Card(
              child: SelectionArea(
                child: ListTile(
                  leading: const Icon(Icons.key),
                  title: Text(name),
                  subtitle: Text(item['ID']!.value),
                  trailing: IconButton(
                    onPressed: () => _showEditModal(context, item['ID']!.value),
                    icon: const Icon(Icons.edit),
                  ),
                ),
              ),
            );
          },
        ),
      ],
    );
  }

  Widget _buildAddButton(BuildContext context) {
    return ElevatedButton(
      onPressed: () => _showEditModal(context, null),
      child: const Text('Create new'),
    );
  }

  void _showEditModal(BuildContext context, String? resourceKey) {
    showDialog(
      context: context,
      builder: (context) {
        return Dialog(
          child: ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: mediumBreakpoint),
            child: Padding(
              padding: C.paddings.defaultPadding,
              child: ResourceForm(
                resourcePath: widget.resourcePath,
                resourceKey: resourceKey,
                resourceSchema: widget.resourceSchema,
                forceCreate: resourceKey == null,
                afterAction: () {
                  Navigator.of(context).pop();
                  ref.invalidate(jsonFetcher(widget.osp()));
                },
              ),
            ),
          ),
        );
      },
    );
  }
}
