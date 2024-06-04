import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

abstract class Module {
  String get description;
  bool get superuserOnly => false;
  NavigationDestination get destination;
  Widget? get summary => null;
  Widget build(BuildContext context);

  String get key => destination.label.toLowerCase().replaceAll(' ', '-');
}

class HomeModule extends Module {
  HomeModule(this.modules);

  final List<Module> modules;

  @override
  String get description => 'Display the summaries of available modules.';

  @override
  NavigationDestination get destination => const NavigationDestination(
        icon: Icon(Icons.home_outlined),
        selectedIcon: Icon(Icons.home),
        label: 'Home',
      );

  @override
  Widget build(BuildContext context) => _buidModuleList(context, modules);
}

Widget _buidModuleList(BuildContext context, List<Module> modules) {
  return Wrap(
    spacing: 8.0,
    runSpacing: 8.0,
    crossAxisAlignment: WrapCrossAlignment.center,
    children: modules
        .where((el) => el is! HomeModule)
        .map((el) => _buildModuleCard(context, el))
        .toList(),
  );
}

Widget _buildModuleCard(BuildContext context, Module module) {
  return Card(
    child: InkWell(
      onTap: () => context.pushNamed('module-${module.key}'),
      child: SizedBox(
        width: 280,
        child: Padding(
          padding: const EdgeInsets.all(16.0),
          child: Column(
            children: [
              Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  module.destination.icon,
                  const SizedBox(width: 8.0),
                  Text(
                    module.destination.label,
                    style: Theme.of(context).textTheme.titleLarge,
                  ),
                ],
              ),
              const SizedBox(height: 16.0),
              Text(module.description),
            ],
          ),
        ),
      ),
    ),
  );
}
