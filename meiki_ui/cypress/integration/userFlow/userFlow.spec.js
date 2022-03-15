/// <reference types="cypress"/>

const baseUrl = "http://localhost:3000/"
const serverUrl = "http://localhost:8080/"

describe("UserFlow Test", () => {
    before(() => {
        const options = {
            method: "POST",
            url: `${serverUrl}delete`,
            failOnStatusCode: false,
            body: {
                username: "shnoo",
                password: "thisisveryunsafe",
            },
        }

        cy.request(options)
    })

    it("Login flow works fully", () => {
        cy.clearCookies()
        cy.visit(baseUrl)

        // shows login page
        cy.get("img[alt='meiki-logo']").should("be.visible")
        cy.get("#Username").should("be.visible")
        cy.get("#Password").should("be.visible")
        cy.get("Button").should("include.text", "Login").and("be.visible")
        cy.get("a").should("include.text", "Create").and("be.visible").click()

        // user creates an account
        cy.get("#Username").type("shnoo")
        cy.get("#Password").type("thisisveryunsafe")
        cy.get("Button")
            .should("include.text", "Create Meiki account")
            .and("be.visible")
            .click()

        // goes to account creation success page
        cy.contains("Your account has successfully been created").should(
            "be.visible"
        )
        cy.get("a").should("be.visible").click()

        // user logs in
        cy.get("#Username").type("shnoo")
        cy.get("#Password").type("thisisveryunsafe")
        cy.get("Button").click()

        // assert it goes to the app
        cy.get("nav").should("be.visible")
        // TODO: adder username in the app
    })
})
