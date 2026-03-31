// Attach lightweight client-side validation to the new sighting form once the DOM is ready.
document.addEventListener("DOMContentLoaded", () => {
    const form = document.getElementById("new-sighting-form");
    if (form) {
        form.addEventListener("submit", (e) => {
            const animal = document.querySelector("#animal").value.trim();
            const location = document.querySelector("#location").value.trim();

            // Block submission when the required fields are blank so the user gets immediate feedback.
            if (!animal || !location) {
                e.preventDefault();
                alert("Please fill in both the animal and location fields.");
            }
        });
    }
});
