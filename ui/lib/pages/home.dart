import 'package:cadanse/cadanse.dart';
import 'package:cadanse/layout.dart';
import 'package:flutter/material.dart';
import 'package:flutter_adaptive_scaffold/flutter_adaptive_scaffold.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../providers/modules.dart';
import '../providers/theme_mode_switch.dart';
import '../services/user_session.dart';

class HomePage extends ConsumerStatefulWidget {
  const HomePage({super.key, this.initialModule});

  final String? initialModule;

  @override
  HomePageState createState() => HomePageState();
}

class HomePageState extends ConsumerState<HomePage> {
  var _selectedIndex = 0;

  @override
  void initState() {
    super.initState();
    final idx = ref
        .read(modulesProvider)
        .all
        .indexWhere((el) => el.key == widget.initialModule);
    if (idx != -1) {
      _selectedIndex = idx;
    }
  }

  @override
  Widget build(BuildContext context) {
    return ref.watch(userSessionProvider).maybeWhen(
          data: (userData) => _buildScaffold(userData!),
          orElse: () => const Center(child: CircularProgressIndicator()),
        );
  }

  Widget _buildScaffold(UserData userData) {
    final modules = ref.watch(modulesProvider);
    return AdaptiveScaffold(
      selectedIndex: _selectedIndex,
      onSelectedIndexChange: (index) =>
          context.pushNamed('module-${modules[index].key}'),
      destinations: modules.all.map((el) => el.destination).toList(),
      body: (context) => Center(
        child: Padding(
          padding: C.paddings.defaultPadding,
          child: Container(
            constraints: const BoxConstraints(maxWidth: mediumBreakpoint),
            child: modules[_selectedIndex].build(context),
          ),
        ),
      ),
      internalAnimations: false,
      appBar: AppBar(
        automaticallyImplyLeading: false,
        title: Row(
          children: [
            const Text('Freon'),
            const Spacer(),
            IconButton(
              onPressed: () =>
                  ref.read(themeModeSwitchProvider.notifier).toggle(),
              icon: ref.watch(themeModeSwitchProvider).toggleIcon,
            ),
            const SizedBox(width: 16.0),
            Text(
              userData.username,
              style: Theme.of(context).textTheme.bodyLarge,
            ),
            const Icon(Icons.person_outlined),
          ],
        ),
      ),
      appBarBreakpoint: Breakpoints.standard,
    );
  }
}
