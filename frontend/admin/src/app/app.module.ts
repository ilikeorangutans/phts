import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { Router, RouterModule, Routes } from "@angular/router";
import { HttpModule } from "@angular/http";

import { CollectionService } from "./services/collection.service";
import { PathService  } from "./services/path.service";

import { AppComponent } from './app.component';

const adminRoutes: Routes = [
  { path: '', component: AppComponent }
]

@NgModule({
  declarations: [
    AppComponent
  ],
  imports: [
    BrowserModule,
    HttpModule
  ],
  providers: [
    CollectionService,
    PathService
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
