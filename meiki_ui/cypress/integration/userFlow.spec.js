/// <reference types="cypress"/>

describe("UserFlow Test", () => {
    beforeEach(() => {
        cy.cleanUsers()
    })

    it("Login flow works fully", () => {
        cy.visit("/login")
        // shows login page
        cy.get("img[alt='meiki-logo']").should("be.visible")
        cy.get("#username").should("be.visible")
        cy.get("#password").should("be.visible")
        cy.get("Button").should("include.text", "Login").and("be.visible")
        cy.get("a[href='/create']").and("be.visible").click()

        // user creates an account
        cy.get("Button")
            .should("include.text", "Create Meiki account")
            .and("be.visible")
        cy.get("#username").type("shnoo")
        cy.get("#password").type("thisisveryunsafe")
        cy.get("#confirmpassword").type("thisisveryunsafe")
        cy.get("Button").click()

        // goes to create success page
        cy.contains("Your account has successfully been created").should(
            "be.visible"
        )
        cy.get("a[href='/login']").should("be.visible").click()

        // user logs in
        cy.get("#username").type("shnoo")
        cy.get("#password").type("thisisveryunsafe")
        cy.get("Button").click()

        // assert it goes to the app
        cy.get("nav").should("be.visible")
        cy.get("[data-cy='profile']").should("contain", "shnoo").click()

        // click logout button
        cy.get("button:contains('Logout')").click()

        // shows login page
        cy.get("img[alt='meiki-logo']").should("be.visible")
        cy.get("#username").should("be.visible")
        cy.get("#password").should("be.visible")
    })

    it("Show password do not match if confirm field does not match", () => {
        cy.visit("/create")

        cy.get("#username").type("shnoo")
        cy.get("#password").type("password")
        cy.get("#confirmpassword").type("doesNotMatch")
        cy.get("Button").click()

        cy.contains("Passwords do not match").should("be.visible")
    })

    it("Show duplicate user error on create", () => {
        cy.createUser("alex", "password")
        cy.visit("/create")

        cy.get("#username").type("alex")
        cy.get("#password").type("password")
        cy.get("#confirmpassword").type("password")
        cy.get("Button").click()

        cy.contains("User already exists").should("be.visible")
    })

    it("Show invalid username error on create", () => {
        cy.visit("/create")

        cy.get("#username").type("alex**")
        cy.get("#password").type("password")
        cy.get("#confirmpassword").type("password")
        cy.get("Button").click()

        cy.contains(
            "Username should not contain any special characters other than '-' and '_'"
        ).should("be.visible")
    })

    it("Show invalid password error on create", () => {
        cy.visit("/create")

        cy.get("#username").type("alex")
        cy.get("#password").type("123")
        cy.get("#confirmpassword").type("123")
        cy.get("Button").click()

        cy.contains("Password should have minimum five characters").should(
            "be.visible"
        )
    })

    it("Error out with unable to connect to server on create", () => {
        cy.simulateServerDown()
        cy.visit("/create")

        cy.get("#username").type("alex")
        cy.get("#password").type("password")
        cy.get("#confirmpassword").type("password")
        cy.get("Button").click()

        cy.contains(
            "An error has occurred while creating the account, unable to connect to server"
        ).should("be.visible")
    })

    it("Show invalid username error on login", () => {
        cy.visit("/login")

        cy.get("#username").type("alex**")
        cy.get("#password").type("password")
        cy.get("Button").click()

        cy.contains(
            "Username should not contain any special characters other than '-' and '_'"
        ).should("be.visible")
    })

    it("Show invalid password error on login", () => {
        cy.visit("/login")

        cy.get("#username").type("alex")
        cy.get("#password").type("123")
        cy.get("Button").click()

        cy.contains("Password should have minimum five characters").should(
            "be.visible"
        )
    })

    it("Show password mismatch on login", () => {
        cy.createUser("alex", "password")
        cy.visit("/login")

        cy.get("#username").type("alex")
        cy.get("#password").type("passwordMismatch")
        cy.get("Button").click()

        cy.contains("Password does not match").should("be.visible")
    })

    it("Error out with unable to connect to server on login", () => {
        cy.simulateServerDown()
        cy.visit("/login")

        cy.get("#username").type("alex")
        cy.get("#password").type("password")
        cy.get("Button").click()

        cy.contains(
            "An error has occurred while logging in, unable to connect to server"
        ).should("be.visible")
    })
})
