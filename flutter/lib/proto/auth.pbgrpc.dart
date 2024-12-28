//
//  Generated code. Do not modify.
//  source: auth.proto
//
// @dart = 2.12

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_final_fields
// ignore_for_file: unnecessary_import, unnecessary_this, unused_import

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:grpc/service_api.dart' as $grpc;
import 'package:protobuf/protobuf.dart' as $pb;

import 'auth.pb.dart' as $0;

export 'auth.pb.dart';

@$pb.GrpcServiceName('auth.AuthService')
class AuthServiceClient extends $grpc.Client {
  static final _$loginWithEmailAndPass = $grpc.ClientMethod<$0.LoginWithEmailAndPassRequest, $0.LoginResponse>(
      '/auth.AuthService/LoginWithEmailAndPass',
      ($0.LoginWithEmailAndPassRequest value) => value.writeToBuffer(),
      ($core.List<$core.int> value) => $0.LoginResponse.fromBuffer(value));
  static final _$register = $grpc.ClientMethod<$0.RegisterRequest, $0.RegisterResponse>(
      '/auth.AuthService/Register',
      ($0.RegisterRequest value) => value.writeToBuffer(),
      ($core.List<$core.int> value) => $0.RegisterResponse.fromBuffer(value));
  static final _$verifyEmail = $grpc.ClientMethod<$0.VerifyEmailRequest, $0.VerifyEmailResponse>(
      '/auth.AuthService/VerifyEmail',
      ($0.VerifyEmailRequest value) => value.writeToBuffer(),
      ($core.List<$core.int> value) => $0.VerifyEmailResponse.fromBuffer(value));

  AuthServiceClient($grpc.ClientChannel channel,
      {$grpc.CallOptions? options,
      $core.Iterable<$grpc.ClientInterceptor>? interceptors})
      : super(channel, options: options,
        interceptors: interceptors);

  $grpc.ResponseFuture<$0.LoginResponse> loginWithEmailAndPass($0.LoginWithEmailAndPassRequest request, {$grpc.CallOptions? options}) {
    return $createUnaryCall(_$loginWithEmailAndPass, request, options: options);
  }

  $grpc.ResponseFuture<$0.RegisterResponse> register($0.RegisterRequest request, {$grpc.CallOptions? options}) {
    return $createUnaryCall(_$register, request, options: options);
  }

  $grpc.ResponseFuture<$0.VerifyEmailResponse> verifyEmail($0.VerifyEmailRequest request, {$grpc.CallOptions? options}) {
    return $createUnaryCall(_$verifyEmail, request, options: options);
  }
}

@$pb.GrpcServiceName('auth.AuthService')
abstract class AuthServiceBase extends $grpc.Service {
  $core.String get $name => 'auth.AuthService';

  AuthServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.LoginWithEmailAndPassRequest, $0.LoginResponse>(
        'LoginWithEmailAndPass',
        loginWithEmailAndPass_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.LoginWithEmailAndPassRequest.fromBuffer(value),
        ($0.LoginResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RegisterRequest, $0.RegisterResponse>(
        'Register',
        register_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.RegisterRequest.fromBuffer(value),
        ($0.RegisterResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.VerifyEmailRequest, $0.VerifyEmailResponse>(
        'VerifyEmail',
        verifyEmail_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.VerifyEmailRequest.fromBuffer(value),
        ($0.VerifyEmailResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.LoginResponse> loginWithEmailAndPass_Pre($grpc.ServiceCall call, $async.Future<$0.LoginWithEmailAndPassRequest> request) async {
    return loginWithEmailAndPass(call, await request);
  }

  $async.Future<$0.RegisterResponse> register_Pre($grpc.ServiceCall call, $async.Future<$0.RegisterRequest> request) async {
    return register(call, await request);
  }

  $async.Future<$0.VerifyEmailResponse> verifyEmail_Pre($grpc.ServiceCall call, $async.Future<$0.VerifyEmailRequest> request) async {
    return verifyEmail(call, await request);
  }

  $async.Future<$0.LoginResponse> loginWithEmailAndPass($grpc.ServiceCall call, $0.LoginWithEmailAndPassRequest request);
  $async.Future<$0.RegisterResponse> register($grpc.ServiceCall call, $0.RegisterRequest request);
  $async.Future<$0.VerifyEmailResponse> verifyEmail($grpc.ServiceCall call, $0.VerifyEmailRequest request);
}
