import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../pages/list.dart';
import 'modules.dart';

class TokensModule extends Module {
  @override
  String get description => 'Manage API tokens.';

  @override
  NavigationDestination get destination => const NavigationDestination(
        icon: Icon(Icons.token_outlined),
        selectedIcon: Icon(Icons.token),
        label: 'Tokens',
      );

  @override
  Widget build(WidgetRef ref, BuildContext context) => const ResourceList(
        resourcePath: '/api/tokens',
        icon: Icon(Icons.key),
        titleKey: 'Name',
      );
}
