import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'modules.dart';

// Handles users, groups and permissions.
class UsersModule extends Module {
  @override
  String get description => 'Users and permissions management.';

  @override
  bool get superuserOnly => true;

  @override
  NavigationDestination get destination => const NavigationDestination(
        icon: Icon(Icons.people_outline),
        selectedIcon: Icon(Icons.people),
        label: 'Users',
      );

  @override
  Widget build(WidgetRef ref, BuildContext context) =>
      const Center(child: Text('Users'));
}
