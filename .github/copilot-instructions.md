---
applyTo: "/frontend"
---

Coding conventions:

# Frontend Coding Conventions for Copilot and Developers

## 🔍 Component Usage

- Before generating new components, **always check the `components/` folder** for reusable components.
- Do **not** create new components if a similar one already exists.

## 🛠 Utility Function Usage

- Before writing new functions, **check the `utils/` folder** for existing utilities.
- Prefer reusing and composing utilities rather than duplicating logic.

## 🎨 Styling Rules

- If an existing component is used, **do not add new styles** unless strictly necessary.
- Reuse the existing styling and follow the project’s design system and CSS conventions.

## 📁 Folder Structure & Naming

- Follow the existing **folder structure** and **naming conventions** strictly.
- Keep components and utilities modular and scoped appropriately.

## ✅ Best Practices

- Avoid code duplication; reuse is preferred.
- Use semantic HTML and accessible markup when applicable.
- Maintain consistent code style and adhere to linting/prettier rules.
