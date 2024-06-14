import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'modules.dart';

class EventsModule extends Module {
  @override
  String get description => 'Explore and manage logs and events.';

  @override
  NavigationDestination get destination => const NavigationDestination(
        icon: Icon(Icons.view_timeline_outlined),
        selectedIcon: Icon(Icons.view_timeline),
        label: 'Events',
      );

  @override
  Widget build(WidgetRef ref, BuildContext context) =>
      const Center(child: Text('Events'));
}
