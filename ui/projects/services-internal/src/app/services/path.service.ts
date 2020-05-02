import { Injectable } from '@angular/core';
import { BasePathService } from 'projects/shared/src/app/services/base-path.service';

@Injectable({
  providedIn: 'root',
})
export class PathService {
  readonly basePath: string;
  constructor(basePathService: BasePathService) {
    this.basePath = basePathService.apiHost;
  }

  servicesPing(): string {
    return [this.basePath, 'services', 'ping'].join('/');
  }
  servicesVersion(): string {
    return [this.basePath, 'services', 'version'].join('/');
  }
}
