describe('Tasks', () => {
  it('Should enable and disable Add button correctly', () => {
    cy.visit('http://localhost:4200');
    cy.get('.taskForm > button').should('not.be.enabled');
    cy.get('.taskInput').type('Task');
    cy.get('.taskForm > button').should('be.enabled');
    cy.get('.taskInput').clear();
    cy.get('.taskForm > button').should('not.be.enabled');
  });

  it('Should filter Pending and Completed tasks correctly', async () => {
    cy.intercept({
      method: 'get',
      url: `/api/tasks/*`,
    }, [])
      .as('getTasks');

    cy.intercept({
      method: 'POST',
      url: '/api/tasks',
    }, { "id": 2, "title": "test", "isCompleted": false })
      .as('createTask');

    cy.intercept({
      method: 'PUT',
      url: '/api/tasks/2'
    }, { "id": 2, "title": "test", "isCompleted": true })
      .as('updateTask');


    cy.visit('http://localhost:4200');
    cy.wait('@getTasks');

    cy.get('.taskInput').type('Task');
    cy.get('.taskForm > button').click();
    cy.wait('@createTask');

    cy.get('.filters > :nth-child(2)').click();
    cy.get('.task').should('be.visible');

    cy.get('.filters > :nth-child(3)').click();
    cy.get('.task').should('not.exist');

    cy.get('.filters > :nth-child(1)').click();
    cy.get('#completed').click();
    cy.wait('@updateTask');

    cy.get('.filters > :nth-child(2)').click();
    cy.get('.task').should('not.exist');

    cy.get('.filters > :nth-child(3)').click();
    cy.get('.task').should('be.visible');

    cy.get('.filters > :nth-child(1)').click();
    cy.get('.task').should('be.visible');
  });
})
