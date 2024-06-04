import 'dart:convert';
import 'dart:html';

import 'package:flutter/foundation.dart';
import 'package:json_annotation/json_annotation.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

import 'freon.dart';

part 'user_session.g.dart';

@riverpod
class UserSession extends _$UserSession {
  @override
  Future<UserData?> build() => fetchData();

  Future<UserData?> fetchData() async {
    try {
      return await freonCall(() async {
        final body = await HttpRequest.getString('$serverUrl/control/user/me');
        return UserData.fromJson(jsonDecode(body));
      });
    } on FreonAuthError catch (_) {
      return null;
    }
  }

  Future<void> login(String username, String password) async {
    await freonCall(() async {
      final data = {
        'username': username,
        'password': password,
      };
      if (kDebugMode) {
        data['next'] =
            window.location.origin + (window.location.pathname ?? '/');
      }
      await HttpRequest.postFormData('$serverUrl/auth/login', data);
    });
    await update((_) => fetchData());
  }

  Future<void> logout() async {
    await freonCall(() async {
      await HttpRequest.request('$serverUrl/auth/logout');
    });
    state = const AsyncValue.data(null);
  }
}

@JsonSerializable(createToJson: false, fieldRename: FieldRename.snake)
class UserData {
  UserData(this.id, this.username, this.isSuperuser);

  final String id;
  final String username;
  final bool isSuperuser;

  factory UserData.fromJson(Map<String, dynamic> json) {
    assert(json['is_superuser']);
    return _$UserDataFromJson(json);
  }
}
