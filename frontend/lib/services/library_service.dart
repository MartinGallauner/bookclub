import '../models/book.dart';

class LibraryService {
  late List<Book> collection;

  LibraryService() {
    collection = [];
  }

  void addBook(Book book) {
    collection.add(book);
  }
}
