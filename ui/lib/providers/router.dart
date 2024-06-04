import 'package:flutter/widgets.dart';
import 'package:go_router/go_router.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

import '../pages/home.dart';
import '../pages/login.dart';
import '../services/user_session.dart';
import 'modules.dart';

part 'router.g.dart';

Page _buildModulePage([String? module]) =>
    NoTransitionPage(child: HomePage(initialModule: module));

@riverpod
GoRouter router(RouterRef ref) {
  GoRouter.optionURLReflectsImperativeAPIs = true;

  final userData = ref.watch(userSessionProvider);
  final modules = ref.watch(modulesProvider).all;
  return GoRouter(routes: [
    GoRoute(
      path: '/',
      name: 'home',
      pageBuilder: (context, state) => _buildModulePage(),
      redirect: (context, state) {
        return userData.maybeWhen(
          data: (value) => value == null ? '/login' : null,
          orElse: () => null,
        );
      },
      routes: modules.map((module) {
        return GoRoute(
          path: module.key,
          name: 'module-${module.key}',
          pageBuilder: (context, state) => _buildModulePage(module.key),
        );
      }).toList(),
    ),
    GoRoute(
      path: '/login',
      name: 'login',
      builder: (context, state) => const LoginPage(),
    ),
  ]);
}
