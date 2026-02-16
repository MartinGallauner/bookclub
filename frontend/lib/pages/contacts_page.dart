import 'package:flutter/material.dart';
import 'package:frontend/models/connection_status.dart';
import 'package:provider/provider.dart';

import '../models/connection.dart';
import '../models/profile.dart';
import '../services/contact_service.dart';

class ContactsPage extends StatelessWidget {
  const ContactsPage({super.key});

  @override
  Widget build(BuildContext context) {
    final contactService = context.watch<ContactService?>();
    final List<Connection> connections = contactService?.connections ?? [];

    return Padding(
      padding: const EdgeInsets.all(16.0),
      child: GridView.builder(
        itemCount: connections.length,
        gridDelegate: const SliverGridDelegateWithMaxCrossAxisExtent(
          maxCrossAxisExtent: 200,
          childAspectRatio: 0.7,
          crossAxisSpacing: 10,
          mainAxisSpacing: 10,
        ),
        itemBuilder: (context, index) {
          //todo forgive me
          Connection connection = connections[index];
          String contactUid = connection.users.firstWhere(
            (id) => id != contactService?.uid,
          );
          return StreamBuilder<Profile?>(
            stream: contactService?.fetchProfile(contactUid),
            builder: (context, snapshot) {
              if (!snapshot.hasData) {
                return Card(child: Center(child: CircularProgressIndicator()));
              }

              Profile contactProfile = snapshot.data!;
              return ContactCard(
                contact: contactProfile,
                connectionStatus: connection.status,
              );
            },
          );
        },
      ),
    );
  }
}

class Contact {
  final String id;
  final String firstName;
  final String lastName;
  final String photoUrl;

  Contact({
    required this.id,
    required this.firstName,
    required this.lastName,
    required this.photoUrl,
  });
}

class ContactCard extends StatelessWidget {
  final Profile contact;
  final ConnectionStatus connectionStatus;

  const ContactCard({
    super.key,
    required this.contact,
    required this.connectionStatus,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      elevation: 4,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Cover image
          Expanded(flex: 3, child: Placeholder()),
          // Book details
          Padding(
            padding: EdgeInsets.all(8.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  contact.displayName,
                  style: TextStyle(fontWeight: FontWeight.bold, fontSize: 14),
                  maxLines: 2,
                  overflow: TextOverflow.ellipsis,
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
