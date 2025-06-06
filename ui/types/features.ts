/**
 * Feature flag interfaces
 */

export interface Features {
  authentication: {
    register: boolean
  }
}

export interface FeaturesResponse {
  features: Features
}
