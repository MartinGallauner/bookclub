import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../models/profile.dart';
import '../services/contact_service.dart';

class ContactsPage extends StatelessWidget {
  const ContactsPage({super.key});

  @override
  Widget build(BuildContext context) {
    final contactState = context.watch<ContactService?>();
    final List<Profile> contacts = contactState?.contacts ?? [];


    return Padding(
      padding: const EdgeInsets.all(16.0),
      child: GridView.builder(
        itemCount: contacts.length,
        gridDelegate: const SliverGridDelegateWithMaxCrossAxisExtent(
          maxCrossAxisExtent: 200,
          childAspectRatio: 0.7,
          crossAxisSpacing: 10,
          mainAxisSpacing: 10,
        ),
        itemBuilder: (context, index) {
          return ContactCard(contact: contacts[index]);
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
  final Contact contact;
  const ContactCard({
    super.key,
    required this.contact,
  });


  @override
  Widget build(BuildContext context) {
    return Card(
      elevation: 4,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Cover image
          Expanded(
            flex: 3,
            child: Placeholder()
          ),
          // Book details
          Padding(
              padding: EdgeInsets.all(8.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                "${contact.firstName} ${contact.lastName}",
                style: TextStyle(
                  fontWeight: FontWeight.bold,
                  fontSize: 14,
                ),
                maxLines: 2,
                overflow: TextOverflow.ellipsis,
              )
            ],
          ))
        ],
      ),

    );

  }

}