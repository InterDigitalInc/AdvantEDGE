module.exports = {
	'env': {
		'browser': true,
		'es6': true
	},
	'extends': 'eslint:recommended',
	'globals': {
		'Atomics': 'readonly',
		'SharedArrayBuffer': 'readonly'
	},
	'parserOptions': {
		'ecmaFeatures': {
			'jsx': true
		},
		'ecmaVersion': 2018,
		'sourceType': 'module'
	},
	'plugins': [
		'react'
	],
	"settings": {
	"react": {
		"createClass": "createReactClass", // Regex for Component Factory to use,
											// default to "createReactClass"
		"pragma": "React",  // Pragma to use, default to "React"
		"version": "15.0", // React version, default to the latest React stable release
		"flowVersion": "0.53" // Flow version
	},
	"propWrapperFunctions": [
			// The names of any function used to wrap propTypes, e.g. `forbidExtraProps`. If this isn't set, any propTypes wrapped in a function will be skipped.
			"forbidExtraProps",
			{"property": "freeze", "object": "Object"},
			{"property": "myFavoriteWrapper"}
		]
	},
	'rules': {
		'indent': [
			'error',
			2
		],
		'linebreak-style': [
			'error',
			'unix'
		],
		'quotes': [
			'error',
			'single'
		],
		'semi': [
			'error',
			'always'
		],
		'extends': [
			'eslint:recommended',
			'plugin:react/recommended'
		],
		"react/jsx-uses-vars": 2,
		"react/jsx-uses-react": "error",
    	// "react/jsx-uses-vars": "error",
	}
};