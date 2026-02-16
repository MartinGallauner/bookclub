import 'package:fake_cloud_firestore/fake_cloud_firestore.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:frontend/services/contact_service.dart';

void main() {
  group('ContactService', () {
    test('Add one contact to empty network.', () async {
      ContactService contactService = ContactService('loggedInUID', FakeFirebaseFirestore());

      var connections = contactService.connections;

      expect(connections, isEmpty);

      await contactService.addContact('newContactUID');
      await Future.delayed(Duration(milliseconds: 100));
      connections = contactService.connections;

      expect(connections.length, equals(1));
      expect(connections[0].uid, equals('loggedInUID_newContactUID'));
      expect(connections[0].requestedBy, equals('loggedInUID'));
      expect(connections[0].acceptedAt, isNull);
      expect(connections[0].rejectedAt, isNull);
    });
  });
}
