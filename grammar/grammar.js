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
      repeat(choice(
        $.prop_path,
        $.key_value_pair
      )),
      '}'
    ),

    task_definition: $ => seq(
      'task',
      field('name', $.string_literal),
      $.task_body
    ),

    task_body: $ => seq(
      '{',
      repeat(choice(
        $.prop_input,
        $.prop_transformer,
        $.key_value_pair
      )),
      '}'
    ),

    sink_definition: $ => seq(
      'sink',
      field('name', $.string_literal),
      $.sink_body
    ),
    sink_body: $ => seq(
      '{',
      repeat(choice(
        $.prop_input,
        $.prop_path,
        $.key_value_pair
      )),
      '}'
    ),

    prop_path: $ => seq(
      'path',
      ':',
      field('path', $.string_literal)
    ),

    prop_input: $ => seq(
      'input',
      ':',
      field('input', $.identifier)
    ),

    prop_transformer: $ => seq(
      'transformer',
      ':',
      field('transformer', $.string_literal)
    ),

    key_value_pair: $ => seq(
      field('key', $.identifier),
      ':',
      field('value', $._value)
    ),

    _value: $ => choice(
      $.string_literal,
      $.identifier,
      $.number,
      $.boolean
    ),

    identifier: $ => /[a-zA-Z_][a-zA-Z0-9_]*/,

    string_literal: $ => token(seq(
      '"',
      repeat(choice(/[^"\\]/, /\\./)),
      '"'
    )),

    number: $ => /\d+/,
    boolean: $ => choice('true', 'false'),

    comment: $ => token(seq('#', /.*/))
  }
});
