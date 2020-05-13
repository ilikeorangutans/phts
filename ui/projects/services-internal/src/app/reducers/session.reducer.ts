import { createReducer, on, Action } from '@ngrx/store';
import { loginSuccess, logout } from '../actions/session.action';

export interface SessionState {
    sessionID: string
    authenticated: boolean
    email: string
};

export const initialState: SessionState = {
    authenticated: false,
    sessionID: "",
    email: "",
};

const _sessionReducer = createReducer(
    initialState,
    on(loginSuccess, (_, { email, sessionID }) => ({ authenticated: true, email: email, sessionID: sessionID })),
    on(logout, _ => initialState),
);

export function sessionReducer(state: SessionState | undefined, action: Action) {
    return _sessionReducer(state, action);
}