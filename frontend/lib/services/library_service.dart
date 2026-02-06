import 'package:flutter/cupertino.dart';

import '../models/book.dart';

class LibraryService extends ChangeNotifier{
  final List<Book> _collection = [];

  List<Book> get books =>  List.unmodifiable(_collection);

  void addBook(Book book) {
    _collection.add(book);
    notifyListeners();
  }

  //TODO remove books
}
