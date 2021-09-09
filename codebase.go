package jsonfilter

// "What the fuck is this??" You might asked.
// Well, I've implemented a Golang version of the compiler
// but it didn't work so well, and I've spent few days working on it,
// yet it still has no progress.
//
// So- yes. Otto (github.com/robertkrimen/otto) is a good JavaScript interpreter.
// The following code from https://github.com/nemtsov/json-mask, man, the author is a nice person.

var js = `
/**
 * util.js
 */

var ObjProto = Object.prototype

function isEmpty (obj) {
  if (obj == null) return true
  if (isArray(obj) ||
     (typeof obj === 'string')) return (obj.length === 0) 
  for (var key in obj) if (has(obj, key)) return false
  return true
}

function isArray (obj) {
  return ObjProto.toString.call(obj) === '[object Array]'
}

function isObject (obj) {
  return (typeof obj === 'function') || (typeof obj === 'object' && !!obj)
}

function has (obj, key) {
  return ObjProto.hasOwnProperty.call(obj, key)
}

/**
 * filter.js
 */

function filter (obj, compiledMask) {
  return isArray(obj)
    ? _arrayProperties(obj, compiledMask)
    : _properties(obj, compiledMask)
}

// wrap array & mask in a temp object;
// extract results from temp at the end
function _arrayProperties (arr, mask) {
  var obj = _properties({ _: arr }, {
    _: {
      type: 'array',
      properties: mask
    }
  })
  return obj && obj._
}

function _properties (obj, mask) {
  var maskedObj, key, value, ret, retKey, typeFunc
  if (!obj || !mask) return obj

  if (isArray(obj)) maskedObj = []
  else if (isObject(obj)) maskedObj = {}

  for (key in mask) {
    if (!has(mask, key)) continue
    value = mask[key]
    ret = undefined
    typeFunc = (value.type === 'object') ? _object : _array
    if (key === '*') {
      ret = _forAll(obj, value.properties, typeFunc)
      for (retKey in ret) {
        if (!has(ret, retKey)) continue
        maskedObj[retKey] = ret[retKey]
      }
    } else {
      ret = typeFunc(obj, key, value.properties)
      if (typeof ret !== 'undefined') maskedObj[key] = ret
    }
  }
  return maskedObj
}

function _forAll (obj, mask, fn) {
  var ret = {}
  var key
  var value
  for (key in obj) {
    if (!has(obj, key)) continue
    value = fn(obj, key, mask)
    if (typeof value !== 'undefined') ret[key] = value
  }
  return ret
}

function _object (obj, key, mask) {
  var value = obj[key]
  if (isArray(value)) return _array(obj, key, mask)
  return mask ? _properties(value, mask) : value
}

function _array (object, key, mask) {
  var ret = []
  var arr = object[key]
  var obj
  var maskedObj
  var i
  var l
  if (!isArray(arr)) return _properties(arr, mask)
  if (isEmpty(arr)) return arr
  for (i = 0, l = arr.length; i < l; i++) {
    obj = arr[i]
    maskedObj = _properties(obj, mask)
    if (typeof maskedObj !== 'undefined') ret.push(maskedObj)
  }
  return ret.length ? ret : undefined
}

/**
 * compiler.js
 */

var TERMINALS = { ',': 1, '/': 2, '(': 3, ')': 4 }

function compile (text) {
  if (!text) return null
  return parse(scan(text))
}

function scan (text) {
  var i = 0
  var len = text.length
  var tokens = []
  var name = ''
  var ch

  function maybePushName () {
    if (!name) return
    tokens.push({ tag: '_n', value: name })
    name = ''
  }

  for (; i < len; i++) {
    ch = text.charAt(i)
    if (TERMINALS[ch]) {
      maybePushName()
      tokens.push({ tag: ch })
    } else {
      name += ch
    }
  }
  maybePushName()

  return tokens
}

function parse (tokens) {
  return _buildTree(tokens, {})
}

function _buildTree (tokens, parent) {
  var props = {}
  var token

  while ((token = tokens.shift())) {
    if (token.tag === '_n') {
      token.type = 'object'
      token.properties = _buildTree(tokens, token)
      if (parent.hasChild) {
        _addToken(token, props)
        return props
      }
    } else if (token.tag === ',') {
      return props
    } else if (token.tag === '(') {
      parent.type = 'array'
      continue
    } else if (token.tag === ')') {
      return props
    } else if (token.tag === '/') {
      parent.hasChild = true
      continue
    }
    _addToken(token, props)
  }

  return props
}

function _addToken (token, props) {
  props[token.value] = { type: token.type }
  if (!isEmpty(token.properties)) {
    props[token.value].properties = token.properties
  }
}

/**
 * index.js
 */

function mask (obj, mask) {
  return filter(obj, compile(mask)) || null
}

var result = JSON.stringify(mask(JSON.parse(inputObj), inputMask))
`
