import 'package:firebase_auth/firebase_auth.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:frontend/models/profile.dart';
import 'package:frontend/services/auth_service.dart';
import 'package:mocktail/mocktail.dart';

class MockFirebaseAuth extends Mock implements FirebaseAuth {}
class MockUserCredential extends Mock implements UserCredential {}
class MockUser extends Mock implements User {}
class FakeAuthProvider extends Fake implements AuthProvider {}

void main() {
  setUpAll(() {
    registerFallbackValue(FakeAuthProvider());
  });

  group('Auth Service', () {
    test('new signup, create new profile', () async {
      final mockAuth = MockFirebaseAuth();
      final authService = AuthService(mockAuth);

      final mockUser = MockUser();
      final mockUserCredential = MockUserCredential();

      when(() => mockUser.uid).thenReturn('test-uid-123');
      when(() => mockUser.email).thenReturn('test@example.com');
      when(() => mockUser.displayName).thenReturn('Test User');
      when(() => mockUser.photoURL).thenReturn('https://example.com/photo.jpg');
      when(() => mockUserCredential.user).thenReturn(mockUser);

      when(() => mockAuth.signInWithPopup(any()))
          .thenAnswer((_) async => mockUserCredential);

      User? user =  (await authService.signInWithGoogle()).user;

      Profile? profile = await authService.fetchProfile('test-uid-123').first;

      expect(profile?.uid, equals('test-uid-123'));
      expect(profile?.email, equals('test@example.com'));
      expect(profile?.displayName, equals('Test User'));
      expect(profile?.photoURL, equals('https://example.com/photo.jpg'));
    });
  });
}
