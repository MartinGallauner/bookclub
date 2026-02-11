import 'package:flutter/material.dart';

class AppNavigationRail extends StatelessWidget {
  const AppNavigationRail({
    super.key,
    required this.selectedIndex,
    required this.onDestinationSelected,
    required this.isExtended,
  });

  final int selectedIndex;
  final ValueChanged<int> onDestinationSelected;
  final bool isExtended;

  @override
  Widget build(BuildContext context) {
    return SafeArea(
      child: NavigationRail(
        extended: isExtended,
        destinations: [
          NavigationRailDestination(
            icon: Icon(Icons.book),
            label: Text('Your Library'),
          ),
          NavigationRailDestination(
            icon: Icon(Icons.search),
            label: Text('Search'),
          ),
          NavigationRailDestination(
            icon: Icon(Icons.person),
            label: Text('Conntacts'),
          ),
        ],
        selectedIndex: selectedIndex,
        onDestinationSelected: onDestinationSelected,
      ),
    );
  }
}