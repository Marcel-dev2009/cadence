import 'package:flutter/material.dart';

class TitleSection extends StatelessWidget {
  const TitleSection({super.key, required this.name});

  final String name;
  @override
  Widget build(BuildContext context) {
    return Padding(
    padding: const EdgeInsets.all(32),
    child: Text(
      name,
      style: const TextStyle(
        fontWeight: FontWeight.w400
      ),
    ),
    );
  }
}
