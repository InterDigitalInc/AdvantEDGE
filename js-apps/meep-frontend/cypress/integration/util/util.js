
// Obtain Cypress data selector formatted string
export function selector(data) {
  return '[data-cy=' + data + ']';
}

// Click on element with provided name (or optionally child with provided text)
export function click(name, text) {
  if (text == null) {
    cy.get(selector(name)).click({force: true});
  } else {
    cy.get(selector(name)).contains(text).click({force: true});
  }
}

// Check element with provided name
export function check(name, check) {
  if (check) {
    cy.get(selector(name)).check({force: true});
  } else {
    cy.get(selector(name)).uncheck({force: true});
  }
}

// Select provided value from drop-down menu
export function type(name, text) {
  // cy.get(selector(name)).clear({force: true}).type(text, {force: true});
  cy.get(selector(name)).clear({force: true}).type('{selectall}{backspace}' + text, {force: true});
}

// Select provided value from drop-down menu
export function select(name, text) {
  cy.get(selector(name)).select(text, {force: true});
}

// Verify that element with provided name includes text as defined in options
export function verify(name, options, value) {
  if (value == null) {
    cy.get(selector(name)).should('exist').and(options);
  } else {
    cy.get(selector(name)).should(options, value);
  }
}

// Verify that element with provided name is in provided state
export function verifyEnabled(name, enabled, timeout) {
  cy.get(selector(name), { timeout: timeout ? timeout : 1000 }).should('exist').and(enabled ? 'not.be.disabled' : 'be.disabled');
}

// Verify that form with provided name is in provided state
export function verifyForm(name, enabled, options, value) {
  cy.get(selector(name)).should('exist').and(enabled ? 'not.be.disabled' : 'be.disabled');
  if (options != null && value != null) {
    cy.get(selector(name)).should(options, value);
  }
}
