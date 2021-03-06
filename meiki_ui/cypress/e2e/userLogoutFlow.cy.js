/// <reference types="cypress"/>

describe("User logout", () => {
    beforeEach(() => {
        cy.login()
    })

    it("User can logout", () => {
        cy.visit("/")

        cy.get("[data-cy='profile']").click()
        cy.contains("Logout").click()
        cy.get("Button").should("include.text", "Login").and("be.visible")
    })

    it("Errors out with unable to connect to server", () => {
        cy.visit("/")
        cy.simulateServerDown("/auth/logout")

        cy.get("[data-cy='profile']").click()
        cy.contains("Logout").click()
        cy.contains("An error has occurred while logging out, unable to connect to server").should(
            "be.visible"
        )
    })
})
