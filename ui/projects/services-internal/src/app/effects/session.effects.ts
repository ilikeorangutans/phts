import { Injectable } from "@angular/core";

import { ofType, Actions, createEffect } from "@ngrx/effects";

import { of } from "rxjs";
import { login, loginSuccess } from '../actions/session.action';
import { exhaustMap } from 'rxjs/operators';

@Injectable()
export class SessionEffects {

    $login = createEffect(() => this.actions$.pipe(
        ofType(login),
        exhaustMap(action => of(loginSuccess({ email: action.username, sessionID: "xxxxxxxxxxxxxxx" })))
    ));

    constructor(
        private actions$: Actions
    ) { }
}