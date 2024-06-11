import 'dart:convert';
import 'dart:html';

import 'package:flutter/foundation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

const serverUrl = kDebugMode ? 'http://localhost:8080' : '';

class FreonError implements Exception {
  FreonError(this.message, [this.xhr]);

  final String message;
  final HttpRequest? xhr;

  @override
  String toString() {
    if (xhr != null) {
      return '$message: ${xhr!.status}: ${xhr!.response}';
    }
    return message;
  }
}

class FreonUnknownError extends FreonError {
  FreonUnknownError(HttpRequest error) : super('Unknown error', error);
}

class FreonAuthError extends FreonError {
  FreonAuthError() : super('Unauthorized');
}

Future<R> freonCall<R>(Future<R> Function() call) async {
  try {
    return await call();
  } on ProgressEvent catch (e) {
    final xhr = e.target as HttpRequest;
    if (xhr.status == 401) {
      throw FreonAuthError();
    }
    throw FreonUnknownError(xhr);
  }
}

Future<dynamic> freonHttpDelete(String path) async {
  return await freonCall(() async {
    final xhr = await HttpRequest.request(
      serverUrl + path,
      method: 'DELETE',
      responseType: 'json',
    );
    // TODO how to cancel xhr if the future is canceled?
    return xhr.response;
  });
}

class ObjectSchemaPath {
  const ObjectSchemaPath(this.path, [this.schemaPath]);

  final String path;
  final String? schemaPath;

  @override
  operator ==(Object other) =>
      other is ObjectSchemaPath &&
      other.path == path &&
      other.schemaPath == schemaPath;

  @override
  int get hashCode => Object.hash(path.hashCode, schemaPath.hashCode);
}

final jsonFetcher =
    FutureProvider.autoDispose.family<dynamic, ObjectSchemaPath>((ref, osp) {
  return freonCall(() async {
    try {
      final xhr = await HttpRequest.request(
        serverUrl + osp.path,
        responseType: 'json',
      );
      ref.onDispose(() => xhr.abort());
      return xhr.response;
    } on ProgressEvent catch (e) {
      final xhr = e.target as HttpRequest;
      if (xhr.status == 404 && osp.schemaPath != null) {
        final xhrSchema = await HttpRequest.request(
          serverUrl + osp.schemaPath!,
          responseType: 'json',
        );
        ref.onDispose(() => xhrSchema.abort());
        return xhrSchema.response;
      }
      rethrow;
    }
  });
});

class ObjectUpload {
  const ObjectUpload(this.path, this.data, {this.method = 'PUT'});

  final String path;
  final String method;
  final Map<String, dynamic> data;

  String get jsonData => jsonEncode(data);
}

Future<dynamic> jsonUpload(
  String path,
  Map<String, dynamic> data, {
  String method = 'PUT',
}) {
  return freonCall(() async {
    final xhr = await HttpRequest.request(
      serverUrl + path,
      method: method,
      responseType: 'json',
      sendData: jsonEncode(data),
    );
    // TODO how to cancel xhr if the future is canceled?
    return xhr.response;
  });
}
