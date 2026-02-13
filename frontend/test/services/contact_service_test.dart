import 'package:fake_cloud_firestore/fake_cloud_firestore.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:frontend/services/contact_service.dart';

void main() {
  group('ContactService', () {
    test('Add one contact to empty network.', () async {
      ContactService contactService = ContactService('user-id-test', FakeFirebaseFirestore());

      var connections = contactService.connections;

      expect(connections, isEmpty);

      await contactService.addContact('new-contact-uid');
      await Future.delayed(Duration(milliseconds: 100));
      connections = contactService.connections;

      expect(connections.length, equals(1));
    });
  });
}
