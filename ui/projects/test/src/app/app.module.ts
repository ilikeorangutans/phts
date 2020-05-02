import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { GRPC_CLIENT_FACTORY, GrpcStandardClientFactory } from '@ngx-grpc/core';
import { GRPC_VERSION_SERVICE_CLIENT_SETTINGS } from '../proto/phts.pbconf';

@NgModule({
  declarations: [
    AppComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule
  ],
  providers: [
    { provide: GRPC_CLIENT_FACTORY, useClass: GrpcStandardClientFactory },
    { provide: GRPC_VERSION_SERVICE_CLIENT_SETTINGS, useValue: { host: 'http://localhost:9999' }},
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
