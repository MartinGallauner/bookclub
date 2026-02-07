import 'dart:async';

import 'package:cloud_firestore/cloud_firestore.dart';
import 'package:flutter/cupertino.dart';

import '../models/book.dart';

class LibraryService extends ChangeNotifier {
  /// The LibraryServices manages the users book collection.
  /// It is a ChangeNotifier because it actually owns the state of the collection.

  String uid = '0';
  List<Book> _collection = [];
  late StreamSubscription<QuerySnapshot> _stream;
  late FirebaseFirestore firebaseFirestore;

  LibraryService(String userId, FirebaseFirestore firebase) {
    uid = userId;
    firebaseFirestore = firebase;
    _stream = firebaseFirestore
        .collection('users')
        .doc(uid)
        .collection('library')
        .snapshots()
        .listen((snapshot) {
          _collection = snapshot.docs.map((doc) {
            var data = doc.data();
            return Book(
              title: data['title'],
              authors: List<String>.from(data['authors']),
              isbn: data['isbn'],
              coverUrl: data['coverUrl'],
              language: data['language'],
              publisher: data['publisher'],
              publishedIn: data['publishedIn'],
              genre: data['genre'],
              description: data['description'],
            );
          }).toList();
          notifyListeners();
        });
  }

  @override
  void dispose() {
    _stream.cancel();
    super.dispose();
  }

  List<Book> get books => List.unmodifiable(_collection);

  Future addBook(Book book) async {
    await firebaseFirestore
        .collection('users')
        .doc(uid)
        .collection('library')
        .doc(book.isbn)
        .set({
          'title': book.title,
          'authors': book.authors,
          'isbn': book.isbn,
          'language': book.language,
          'publisher': book.publisher,
          'publishedIn': book.publishedIn,
          'genre': book.genre,
          'description': book.description,
          'coverUrl': book.coverUrl,
          'addedAt': FieldValue.serverTimestamp(),
        });
  }
}
