module.exports = grammar({
  name: 'flow',

  extras: $ => [
    $.comment,
    /\s/
  ],

  rules: {
    source_file: $ => repeat($._definition),

    _definition: $ => choice(
      $.source_definition,
      $.task_definition,
      $.sink_definition
    ),

    source_definition: $ => seq(
      'source',
      field('name', $.string_literal),
      $.source_body
    ),
    source_body: $ => seq(
      '{',
      repeat($.key_value_pair),
      $.source_body_path,
      repeat($.key_value_pair),
      '}'
    ),
    source_body_path: $ => seq(
      'path',
      ':',
      field('path', $.string_literal),
    ),

    task_definition: $ => seq(
      'task',
      field('name', $.string_literal),
      $.task_body
    ),
    task_body: $ => seq(
      '{',
      repeat($.key_value_pair),
      $.task_body_input,
      repeat($.key_value_pair),
      $.task_body_transformer,
      repeat($.key_value_pair),
      '}'
    ),
    task_body_input: $ => seq(
      'input',
      ':',
      field('input', $.identifier),
    ),
    task_body_transformer: $ => seq(
      'transformer',
      ':',
      field('transformer', $.string_literal),
    ),

    sink_definition: $ => seq(
      'sink',
      field('name', $.string_literal),
      $.sink_body
    ),
    sink_body: $ => seq(
      '{',
      repeat($.key_value_pair),
      $.sink_body_input,
      repeat($.key_value_pair),
      $.sink_body_path,
      repeat($.key_value_pair),
      '}'
    ),
    sink_body_input: $ => seq(
      'input',
      ':',
      field('input', $.identifier),
    ),
    sink_body_path: $ => seq(
      'path',
      ':',
      field('path', $.string_literal),
    ),


    key_value_pair: $ => seq(
      field('key', $.identifier),
      ':',
      field('value', $._value)
    ),

    _value: $ => choice(
      $.string_literal,
      $.identifier
    ),

    // --- terminal symbols ---

    identifier: $ => /[a-zA-Z_][a-zA-Z0-9_]*/,

    string_literal: $ => token(seq(
      '"',
      repeat(choice(/[^"\\]/, /\\./)),
      '"'
    )),

    comment: $ => token(seq('#', /.*/))
  }
});
