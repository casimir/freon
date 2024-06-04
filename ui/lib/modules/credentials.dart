import 'package:flutter/material.dart';

import '../pages/forms.dart';
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
  Widget build(BuildContext context) =>
      const ResourceForm(resourcePath: '/wallabag/credentials');
}
