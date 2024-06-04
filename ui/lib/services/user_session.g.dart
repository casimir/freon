// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'user_session.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

UserData _$UserDataFromJson(Map<String, dynamic> json) => UserData(
      json['id'] as String,
      json['username'] as String,
      json['is_superuser'] as bool,
    );

// **************************************************************************
// RiverpodGenerator
// **************************************************************************

String _$userSessionHash() => r'fce8dab9a1097aa83971656c467feb51e7949ef2';

/// See also [UserSession].
@ProviderFor(UserSession)
final userSessionProvider =
    AutoDisposeAsyncNotifierProvider<UserSession, UserData?>.internal(
  UserSession.new,
  name: r'userSessionProvider',
  debugGetCreateSourceHash:
      const bool.fromEnvironment('dart.vm.product') ? null : _$userSessionHash,
  dependencies: null,
  allTransitiveDependencies: null,
);

typedef _$UserSession = AutoDisposeAsyncNotifier<UserData?>;
// ignore_for_file: type=lint
// ignore_for_file: subtype_of_sealed_class, invalid_use_of_internal_member, invalid_use_of_visible_for_testing_member
