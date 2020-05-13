import { createAction, props } from '@ngrx/store';

export const login = createAction('[Session] login', props<{ username: string, password: string }>());
export const loginFailure = createAction('[Session] login failure', props<{ message: string }>());
export const loginSuccess = createAction('[Session] login success', props<{ email: string, sessionID: string }>());
export const logout = createAction('[Session] logout');