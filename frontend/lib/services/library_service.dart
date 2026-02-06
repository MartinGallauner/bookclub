import 'package:flutter/cupertino.dart';

import '../models/book.dart';

class LibraryService extends ChangeNotifier{
  /// The LibraryServices manages the users book collection.
  /// It is a ChangeNotifier because it actually owns the state of the collection.
  final List<Book> _collection = [];

  List<Book> get books =>  List.unmodifiable(_collection);

  void addBook(Book book) {
    _collection.add(book);
    print('library service: added book ${book.title}');
    notifyListeners();
  }

  //TODO remove books
}
