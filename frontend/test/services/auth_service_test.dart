
import 'package:bookclub_api/bookclub_api.dart';
import 'package:dio/dio.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:frontend/services/api_client.dart';
import 'package:frontend/services/auth_service.dart';
import 'package:google_sign_in/google_sign_in.dart';
import 'package:mocktail/mocktail.dart';

class MockGoogleSignIn extends Mock implements GoogleSignIn {}
class MockApiClient extends Mock implements ApiClient {}
class MockAuthAPI extends Mock implements AuthApi {}
class MockGoogleSignInAccount extends Mock implements GoogleSignInAccount {}
class MockGoogleSignInAuthentication extends Mock implements GoogleSignInAuthentication {}



void main() {
  setUpAll(() {
    registerFallbackValue(GoogleLoginRequest((b) => b.idToken = ''));
  });
  group('Auth Service', () {
    test('signInWithGoogle returns UserProfile on success', () async {

      final mockAuthApi = MockAuthAPI();
      when(() => mockAuthApi.loginWithGoogle(
        googleLoginRequest: any(named: 'googleLoginRequest'),
      )).thenAnswer((_) async => Response(
        data: AuthResponse((b) => b
          ..token = 'mock-google-token'
          ..tokenType = AuthResponseTokenTypeEnum.bearer
          ..expiresIn = 3600
          ..user = (UserProfileBuilder()
            ..id = 'test-user-id'
            ..email = 'test@example.com'
            ..displayName = 'Test User'
          )
        ),
        requestOptions: RequestOptions(),
        statusCode: 200,
      ));




      final mockApiClient = MockApiClient();
      when(() => mockApiClient.authApi).thenReturn(mockAuthApi);

      final mockGoogleSignInAuthentication = MockGoogleSignInAuthentication();
      when(() => mockGoogleSignInAuthentication.idToken).thenReturn("fake-id-token");

      final mockGoogleSignInAccount = MockGoogleSignInAccount();
      when(() => mockGoogleSignInAccount.authentication).thenReturn(mockGoogleSignInAuthentication);

      final mockGoogleSignIn = MockGoogleSignIn();
      when(() => mockGoogleSignIn.authenticate()).thenAnswer((_) async => mockGoogleSignInAccount);

      final authService = AuthService(mockApiClient, mockGoogleSignIn);

      var result = await authService.signInWithGoogle();

      expect(result.id, equals('test-user-id'));
      expect(result.email, equals('test@example.com'));

    });

    test('signInWithGoogle throws if Google sign-in returns null', () async {

    });

  });
}
