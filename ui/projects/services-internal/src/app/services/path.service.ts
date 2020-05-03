import { Injectable } from '@angular/core';
import { BasePathService } from 'projects/shared/src/app/services/base-path.service';

@Injectable({
  providedIn: 'root',
})
export class PathService {
  readonly basePath: string;
  readonly services: string;
  readonly servicesInternal: string;

  constructor(basePathService: BasePathService) {
    this.basePath = basePathService.apiHost;
    this.services = [this.basePath, 'services'].join('/');
    this.servicesInternal = [this.services, 'internal'].join('/');
  }

  sessionCreate(): string {
    return [this.servicesInternal, 'sessions', 'create'].join('/');
  }

  sessionDestroy(): string {
    return [this.servicesInternal, 'sessions', 'destroy'].join('/');
  }

  servicesPing(): string {
    return [this.basePath, 'services', 'ping'].join('/');
  }

  servicesVersion(): string {
    return [this.basePath, 'services', 'version'].join('/');
  }
}
