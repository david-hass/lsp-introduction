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
      field('body', $.block_body),
    ),

    task_definition: $ => seq(
      'task',
      field('name', $.string_literal),
      field('body', $.block_body),
    ),

    sink_definition: $ => seq(
      'sink',
      field('name', $.string_literal),
      field('body', $.block_body),
    ),

    block_body: $ => seq(
      '{',
      repeat($.key_value_pair),
      '}'
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
