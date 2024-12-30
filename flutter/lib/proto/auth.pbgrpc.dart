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
  static final _$resendVerificationEmail = $grpc.ClientMethod<$0.ResendVerificationEmailRequest, $0.ResendVerificationEmailResponse>(
      '/auth.AuthService/ResendVerificationEmail',
      ($0.ResendVerificationEmailRequest value) => value.writeToBuffer(),
      ($core.List<$core.int> value) => $0.ResendVerificationEmailResponse.fromBuffer(value));
  static final _$changePassword = $grpc.ClientMethod<$0.ChangePasswordRequest, $0.ChangePasswordResponse>(
      '/auth.AuthService/ChangePassword',
      ($0.ChangePasswordRequest value) => value.writeToBuffer(),
      ($core.List<$core.int> value) => $0.ChangePasswordResponse.fromBuffer(value));
  static final _$changeEmail = $grpc.ClientMethod<$0.ChangeEmailRequest, $0.ChangeEmailResponse>(
      '/auth.AuthService/ChangeEmail',
      ($0.ChangeEmailRequest value) => value.writeToBuffer(),
      ($core.List<$core.int> value) => $0.ChangeEmailResponse.fromBuffer(value));

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

  $grpc.ResponseFuture<$0.ResendVerificationEmailResponse> resendVerificationEmail($0.ResendVerificationEmailRequest request, {$grpc.CallOptions? options}) {
    return $createUnaryCall(_$resendVerificationEmail, request, options: options);
  }

  $grpc.ResponseFuture<$0.ChangePasswordResponse> changePassword($0.ChangePasswordRequest request, {$grpc.CallOptions? options}) {
    return $createUnaryCall(_$changePassword, request, options: options);
  }

  $grpc.ResponseFuture<$0.ChangeEmailResponse> changeEmail($0.ChangeEmailRequest request, {$grpc.CallOptions? options}) {
    return $createUnaryCall(_$changeEmail, request, options: options);
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
    $addMethod($grpc.ServiceMethod<$0.ResendVerificationEmailRequest, $0.ResendVerificationEmailResponse>(
        'ResendVerificationEmail',
        resendVerificationEmail_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ResendVerificationEmailRequest.fromBuffer(value),
        ($0.ResendVerificationEmailResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ChangePasswordRequest, $0.ChangePasswordResponse>(
        'ChangePassword',
        changePassword_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ChangePasswordRequest.fromBuffer(value),
        ($0.ChangePasswordResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ChangeEmailRequest, $0.ChangeEmailResponse>(
        'ChangeEmail',
        changeEmail_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ChangeEmailRequest.fromBuffer(value),
        ($0.ChangeEmailResponse value) => value.writeToBuffer()));
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

  $async.Future<$0.ResendVerificationEmailResponse> resendVerificationEmail_Pre($grpc.ServiceCall call, $async.Future<$0.ResendVerificationEmailRequest> request) async {
    return resendVerificationEmail(call, await request);
  }

  $async.Future<$0.ChangePasswordResponse> changePassword_Pre($grpc.ServiceCall call, $async.Future<$0.ChangePasswordRequest> request) async {
    return changePassword(call, await request);
  }

  $async.Future<$0.ChangeEmailResponse> changeEmail_Pre($grpc.ServiceCall call, $async.Future<$0.ChangeEmailRequest> request) async {
    return changeEmail(call, await request);
  }

  $async.Future<$0.LoginResponse> loginWithEmailAndPass($grpc.ServiceCall call, $0.LoginWithEmailAndPassRequest request);
  $async.Future<$0.RegisterResponse> register($grpc.ServiceCall call, $0.RegisterRequest request);
  $async.Future<$0.VerifyEmailResponse> verifyEmail($grpc.ServiceCall call, $0.VerifyEmailRequest request);
  $async.Future<$0.ResendVerificationEmailResponse> resendVerificationEmail($grpc.ServiceCall call, $0.ResendVerificationEmailRequest request);
  $async.Future<$0.ChangePasswordResponse> changePassword($grpc.ServiceCall call, $0.ChangePasswordRequest request);
  $async.Future<$0.ChangeEmailResponse> changeEmail($grpc.ServiceCall call, $0.ChangeEmailRequest request);
}
