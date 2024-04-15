window.onload = () => {
    initializeLogout()
    greatModalWindow()
    getFormValue()
}

function initializeLogout() {
    const logoutButton = document.querySelector(".header-button")
    logoutButton.addEventListener('click', async () => {
        const response = await fetch('/api/logout')
        if (response.ok) {
            window.location.href = "/login"
        }
    })
}

function greatModalWindow() {
    const openButtons = document.querySelectorAll(".link-change")
    const modalWindow = document.querySelector(".modal")
    const closeButtons = document.querySelectorAll(".modal-cross")

    openButtons.forEach(openButton => {
        openButton.addEventListener('click', () => {
            modalWindow.classList.add("active")
        })
    })

    closeButtons.forEach(closeButton => {
        closeButton.addEventListener('click', () => {
            modalWindow.classList.remove("active")
        })
    })
}

function getFormValue() {
    let Data = {
        customerId: "",
        status: "",
    }

    const openModalBtns = document.querySelectorAll('.open-modal-btn')
    if (openModalBtns.length > 0) {
        openModalBtns.forEach(function (btn) {
            btn.addEventListener('click', function () {
                const customerId = btn.getAttribute('data-customer-id')
                const statusValue = form.querySelector('select[name="status"]').value
                Data.customerId = customerId
                Data.status = statusValue
            })
        })
    }
    const form = document.querySelector(".modal-form")
    const statusSelect = form.querySelector('select[name="status"]')
    statusSelect.addEventListener('change', function () {
        let status = statusSelect.value;
        Data.status = status;
    })
    const updateButton = document.querySelector(".modal-form__button")
    updateButton.addEventListener('click', async () => {
        try {
            const response = await fetch('/api/change', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json;charset=utf-8'
                },
                body: JSON.stringify(Data)
            })

            if (response.ok) {
                console.log("Status change was successful")
            } else {
                const errorMessage = await response.text()
                console.error(errorMessage)
            }
        } catch (error) {
            console.error('Error:', error)
        }
    })
}
