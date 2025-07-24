# 🔑 GitHub Actions Tokens Setup Guide

This guide explains how to obtain and configure all the tokens needed for your TOPAY-Z512 GitHub Actions workflows.

## 📋 Required Tokens

| Token | Purpose | Status | Required For |
|-------|---------|--------|--------------|
| `GITHUB_TOKEN` | Repository operations | ✅ Auto-provided | Releases, Pages deployment |
| `NPM_TOKEN` | npm package publishing | ⚠️ Manual setup | JavaScript package publishing |
| `CARGO_REGISTRY_TOKEN` | crates.io publishing | ⚠️ Manual setup | Rust crate publishing |
| `SONAR_TOKEN` | Code quality analysis | ⚠️ Manual setup | SonarCloud integration |

## 🚀 Step-by-Step Setup

### 1. NPM_TOKEN (JavaScript Package Publishing)

#### Get the Token

1. **Create npm account**: Visit [npmjs.com](https://www.npmjs.com) and sign up
2. **Verify email**: Check your email and verify your account
3. **Generate token**:

   ```bash
   npm login
   npm token create --read-and-publish
   ```

   **Or via Web UI**:
   - Go to npmjs.com → Profile → Access Tokens
   - Click "Generate New Token"
   - Select "Automation" type
   - Copy the token (starts with `npm_`)

#### Add to GitHub

1. Go to your repository: `https://github.com/TOPAY-FOUNDATION/TOPAY_Z512`
2. Click **Settings** → **Secrets and variables** → **Actions**
3. Click **New repository secret**
4. Name: `NPM_TOKEN`
5. Value: Your npm token
6. Click **Add secret**

### 2. CARGO_REGISTRY_TOKEN (Rust Crate Publishing)

#### Obtaining the Cargo Registry Token

1. **Create crates.io account**: Visit [crates.io](https://crates.io)
2. **Sign in with GitHub**: Use your GitHub account
3. **Generate API token**:
   - Go to Account Settings → API Tokens
   - Click "New Token"
   - Name: "GitHub Actions TOPAY_Z512"
   - Copy the generated token

#### Adding Cargo Registry Token to GitHub

1. Repository → **Settings** → **Secrets and variables** → **Actions**
2. Click **New repository secret**
3. Name: `CARGO_REGISTRY_TOKEN`
4. Value: Your crates.io token
5. Click **Add secret**

### 3. SONAR_TOKEN (Code Quality Analysis)

#### Obtaining the SonarCloud Token

1. **Create SonarCloud account**: Visit [sonarcloud.io](https://sonarcloud.io)
2. **Sign up with GitHub**: Connect your GitHub account
3. **Import project**:
   - Click "+" → "Analyze new project"
   - Select your `TOPAY_Z512` repository
   - Choose "With GitHub Actions"
4. **Generate token**:
   - Go to My Account → Security → Generate Tokens
   - Name: "GitHub Actions"
   - Copy the token

#### Adding SonarCloud Token to GitHub

1. Repository → **Settings** → **Secrets and variables** → **Actions**
2. Click **New repository secret**
3. Name: `SONAR_TOKEN`
4. Value: Your SonarCloud token
5. Click **Add secret**

## ✅ Verification

### Test Your Setup

1. Go to your repository → **Actions** tab
2. Click **Test Tokens Configuration** workflow
3. Click **Run workflow** → **Run workflow**
4. Check the results to see which tokens are properly configured

### Expected Results

```output
✅ GITHUB_TOKEN is available
✅ NPM_TOKEN is available  
✅ CARGO_REGISTRY_TOKEN is available
✅ SONAR_TOKEN is available
```

## 🔧 Package Configuration

Your packages are now properly configured with the correct repository URLs:

### JavaScript Package (`js/package.json`)

- **Name**: `@topayfoundation/topayz512`
- **Repository**: `https://github.com/TOPAY-FOUNDATION/TOPAY_Z512`
- **Registry**: npm (npmjs.com)

### Rust Crate (`rust/Cargo.toml`)

- **Name**: `topayz512`
- **Repository**: `https://github.com/TOPAY-FOUNDATION/TOPAY_Z512`
- **Registry**: crates.io

## 🚨 Security Best Practices

1. **Never commit tokens**: Tokens should only be stored in GitHub Secrets
2. **Use minimal permissions**: Only grant necessary permissions to tokens
3. **Rotate tokens regularly**: Update tokens periodically for security
4. **Monitor usage**: Check token usage in respective platforms

## 🔄 Workflow Triggers

Your workflows will automatically run when:

- **CI/CD** (`ci.yml`): On push/PR to main branch
- **Release** (`release.yml`): When you create version tags (`v*`, `js/v*`, `go/v*`, `rust/v*`)
- **Quality** (`quality.yml`): On push/PR to main branch
- **Security** (`security.yml`): Daily and on push/PR
- **Nightly** (`nightly.yml`): Daily at midnight UTC

## 📞 Support

If you encounter issues:

1. Check the **Actions** tab for detailed error logs
2. Verify all secrets are properly configured
3. Ensure package configurations are correct
4. Run the token test workflow to diagnose issues

---

**Next Steps**: Once all tokens are configured, your workflows will automatically handle building, testing, and publishing your packages! 🎉
