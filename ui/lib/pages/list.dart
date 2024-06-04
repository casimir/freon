import 'package:cadanse/components/widgets/cards.dart';
import 'package:cadanse/components/widgets/error.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../services/freon.dart';
import 'forms.dart';

class ResourceList extends ConsumerStatefulWidget {
  const ResourceList({super.key, required this.resourcePath});

  final String resourcePath;

  @override
  ConsumerState<ResourceList> createState() => _ResourceListState();
}

class _ResourceListState extends ConsumerState<ResourceList> {
  @override
  Widget build(BuildContext context) {
    final url = '/control${widget.resourcePath}';
    return ref.watch(jsonFetcher(url)).when(
          data: (data) => _buildList(data),
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
    return ListView.builder(
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
            ),
          ),
        );
      },
    );
  }
}
