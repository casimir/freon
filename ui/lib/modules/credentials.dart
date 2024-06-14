import 'package:cadanse/cadanse.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../pages/forms.dart';
import '../services/freon.dart';
import 'modules.dart';

class CredentialsModule extends Module {
  @override
  String get description => 'Manage credentials for wallabag.';

  @override
  NavigationDestination get destination => const NavigationDestination(
        icon: Icon(Icons.leak_add_outlined),
        selectedIcon: Icon(Icons.leak_add),
        label: 'Credentials',
      );

  @override
  Widget build(WidgetRef ref, BuildContext context) {
    return Column(
      children: [
        Row(
          children: [
            const Spacer(),
            const Text(
              'Session state',
              style: TextStyle(fontWeight: FontWeight.bold),
            ),
            C.spacers.horizontalContent,
            ref.watch(wallabagSessionCheckProvider).when(
                  data: (result) {
                    return result == null
                        ? const Icon(Icons.check_circle, color: Colors.green)
                        : const Icon(Icons.error, color: Colors.orange);
                  },
                  error: (e, st) => Tooltip(
                    message: e.toString(),
                    child: const Icon(Icons.error, color: Colors.red),
                  ),
                  loading: () => const CircularProgressIndicator(),
                ),
            C.spacers.horizontalContent,
            IconButton(
              onPressed: () => ref.invalidate(wallabagSessionCheckProvider),
              icon: const Icon(Icons.refresh),
            )
          ],
        ),
        C.spacers.verticalContent,
        const ResourceForm(resourcePath: '/wallabag/credentials'),
      ],
    );
  }
}
