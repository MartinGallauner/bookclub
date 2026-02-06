import '../models/book.dart';

class LibraryService {
  final List<Book> _collection = [];

  List<Book> get books =>  List.unmodifiable(_collection);

  void addBook(Book book) {
    _collection.add(book);
  }
}
