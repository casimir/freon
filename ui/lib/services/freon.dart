import 'dart:html';

import 'package:flutter/foundation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

const serverUrl = kDebugMode ? 'http://localhost:8080' : '';

class FreonError implements Exception {
  FreonError(this.message, [this.error]);

  final String message;
  final Object? error;

  @override
  String toString() {
    if (error != null) {
      return '$message: $error';
    }
    return message;
  }
}

class FreonUnknownError extends FreonError {
  FreonUnknownError(Object error) : super('Unknown error', error);
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
    throw FreonUnknownError(e);
  }
}

final jsonFetcher =
    FutureProvider.autoDispose.family<dynamic, String>((ref, path) {
  return freonCall(() async {
    final url = serverUrl + path;
    final xhr = await HttpRequest.request(url, responseType: 'json');
    ref.onDispose(() {
      xhr.abort();
    });
    return xhr.response;
  });
});

Future<HttpRequest> postForm(String path, Map<String, String> data) async {
  return await freonCall(
      () async => HttpRequest.postFormData('$serverUrl$path', data));
}
