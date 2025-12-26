import 'package:flutter/material.dart';

class NetworkPage extends StatelessWidget {
  const NetworkPage({super.key});

  @override
  Widget build(BuildContext context) {
    // Mock data for testing
    final List<Contact> contacts = [
      Contact(id: 'a3f2b8c1-4d5e-4a1b-9c2d-1e3f4a5b6c7d', firstName: 'Edward', lastName: 'Teach', photoUrl: ''),
      Contact(id: 'b7d3c9e2-5f6a-4b2c-8d3e-2f4a5b6c7d8e', firstName: 'Anne', lastName: 'Bonny', photoUrl: ''),
      Contact(id: 'c1e4d0f3-6a7b-4c3d-9e4f-3a5b6c7d8e9f', firstName: 'Mary', lastName: 'Read', photoUrl: ''),
      Contact(id: 'd5f8e1a4-7b8c-4d4e-0f5a-4b6c7d8e9f0a', firstName: 'William', lastName: 'Kidd', photoUrl: ''),
      Contact(id: 'e9a2f5b8-8c9d-4e5f-1a6b-5c7d8e9f0a1b', firstName: 'Bartholomew', lastName: 'Roberts', photoUrl: ''),
      Contact(id: 'f0b6a3c9-9d0e-4f6a-2b7c-6d8e9f0a1b2c', firstName: 'Henry', lastName: 'Morgan', photoUrl: ''),
      Contact(id: '1c7d4e0a-0e1f-4a7b-3c8d-7e9f0a1b2c3d', firstName: 'Grace', lastName: 'O\'Malley', photoUrl: ''),
      Contact(id: '2d8e5f1b-1f2a-4b8c-4d9e-8f0a1b2c3d4e', firstName: 'Jack', lastName: 'Rackham', photoUrl: ''),
      Contact(id: '3e9f6a2c-2a3b-4c9d-5e0f-9a1b2c3d4e5f', firstName: 'Charles', lastName: 'Vane', photoUrl: ''),
      Contact(id: '4f0a7b3d-3b4c-4d0e-6f1a-0b2c3d4e5f6a', firstName: 'Stede', lastName: 'Bonnet', photoUrl: ''),
      Contact(id: '5a1b8c4e-4c5d-4e1f-7a2b-1c3d4e5f6a7b', firstName: 'Henry', lastName: 'Every', photoUrl: ''),
      Contact(id: '6b2c9d5f-5d6e-4f2a-8b3c-2d4e5f6a7b8c', firstName: 'Francis', lastName: 'Drake', photoUrl: ''),
    ];

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