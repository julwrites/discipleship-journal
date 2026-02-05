import '@testing-library/jest-dom';

// Mock localStorage for tests
const localStorageMock = {
    store: {} as Record<string, string>,
    clear() {
        this.store = {};
    },
    getItem(key: string) {
        return this.store[key] || null;
    },
    setItem(key: string, value: string) {
        this.store[key] = value.toString();
    },
    removeItem(key: string) {
        delete this.store[key];
    },
    get length() {
        return Object.keys(this.store).length;
    },
    key(index: number) {
        return Object.keys(this.store)[index] || null;
    }
};

Object.defineProperty(window, 'localStorage', {
    value: localStorageMock,
});

// Mock ResizeObserver
global.ResizeObserver = class ResizeObserver {
  observe() {}
  unobserve() {}
  disconnect() {}
};

// Mock scrollIntoView
Element.prototype.scrollIntoView = () => {};
