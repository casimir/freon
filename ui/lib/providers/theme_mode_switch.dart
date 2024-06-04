import 'package:flutter/material.dart';
import 'package:flutter/scheduler.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

part 'theme_mode_switch.g.dart';

class ThemeModeData {
  const ThemeModeData(this.mode);

  factory ThemeModeData.fromBrightness(Brightness brightness) => ThemeModeData(
      brightness == Brightness.light ? ThemeMode.light : ThemeMode.dark);

  final ThemeMode mode;

  Icon get toggleIcon => mode == ThemeMode.light
      ? const Icon(Icons.dark_mode_outlined)
      : const Icon(Icons.light_mode_outlined);
}

@riverpod
class ThemeModeSwitch extends _$ThemeModeSwitch {
  @override
  ThemeModeData build() => ThemeModeData.fromBrightness(
      SchedulerBinding.instance.platformDispatcher.platformBrightness);

  Future<void> toggle() async {
    state = ThemeModeData(
        state.mode == ThemeMode.light ? ThemeMode.dark : ThemeMode.light);
  }
}
