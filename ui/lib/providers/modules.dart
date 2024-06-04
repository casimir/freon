import 'package:riverpod_annotation/riverpod_annotation.dart';

import '../modules/credentials.dart';
import '../modules/events.dart';
import '../modules/modules.dart';
import '../modules/tokens.dart';
import '../modules/users.dart';
import '../services/user_session.dart';

part 'modules.g.dart';

final availableModules = [
  CredentialsModule(),
  TokensModule(),
  EventsModule(),
  UsersModule(),
];

class ModuleIndex {
  ModuleIndex(this.modules) : home = HomeModule(modules);

  final Module home;
  final List<Module> modules;

  List<Module> get all => [home, ...modules];
  Module operator [](int index) => all[index];
}

@riverpod
ModuleIndex modules(ModulesRef ref) {
  return ref.watch(userSessionProvider).maybeWhen(
        data: (userData) {
          if (userData == null) {
            return ModuleIndex([]);
          }
          final modules = availableModules
              .where(
                (el) => (!el.superuserOnly ||
                    el.superuserOnly == userData.isSuperuser),
              )
              .toList();
          return ModuleIndex(modules);
        },
        orElse: () => ModuleIndex([]),
      );
}
