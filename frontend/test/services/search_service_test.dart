import 'package:fake_cloud_firestore/fake_cloud_firestore.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:frontend/models/connection.dart';
import 'package:frontend/models/connection_status.dart';
import 'package:frontend/models/profile.dart';
import 'package:frontend/services/auth_service.dart';
import 'package:frontend/services/contact_service.dart';
import 'package:frontend/services/search_service.dart';
import 'package:mocktail/mocktail.dart';

class MockContactService extends Mock implements ContactService {}

class MockAuthService extends Mock implements AuthService {}

void main() {
  group('Search Services Tests', () {
    test('find book in one of two contacts.', () async {
      final mockContactService = MockContactService();
      final Connection connection1 = Connection(
        requestedBy: 'testUser',
        requestedAt: DateTime(2025),
        users: ['testUser', 'userA'],
        status: ConnectionStatus.accepted,
      );
      final Connection connection2 = Connection(
        requestedBy: 'testUser',
        requestedAt: DateTime(2025),
        users: ['testUser', 'userB'],
        status: ConnectionStatus.accepted,
      );
      final connections = [connection1, connection2];
      when(() => mockContactService.connections).thenReturn(connections);
      when(() => mockContactService.uid).thenReturn('testUser');

      final AuthService mockAuthService = MockAuthService();
      when(() => mockAuthService.fetchProfile('userA')).thenAnswer(
        (_) async => Profile(
          uid: 'userA',
          email: 'userA@test.com',
          displayName: 'Herbert Test',
          photoURL: 'photoURL',
          updatedAt: DateTime(2024),
          createdAt: DateTime(2023),
        ),
      );
      when(() => mockAuthService.fetchProfile('userB')).thenAnswer(
        (_) async => Profile(
          uid: 'userB',
          email: 'userB@test.com',
          displayName: 'Thomas Test',
          photoURL: 'photoURL',
          updatedAt: DateTime(2024),
          createdAt: DateTime(2023),
        ),
      );

      final mockFirestore = FakeFirebaseFirestore();
      await mockFirestore
          .collection('users')
          .doc('userB')
          .collection('library')
          .doc('0451524934')
          .set({
            'title': 'Nineteen Eighty-four',
            'authors': ['George Orwell'],
            'isbn': '0451524934',
            'coverUrl': "",
            'language': 'en',
          });

      final service = SearchService(
        mockContactService,
        mockAuthService,
        mockFirestore,
      );

      List<Profile> profiles = await service.searchByISBN('0451524934');

      expect(profiles[0].uid, 'userB');
    });
  });
}
